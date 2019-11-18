// This is a modified version of bgs.c from https://github.com/Gottox/bgs (also MIT licensed).
// Credits are given in the LICENSE file.

#include <math.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>

#include <Imlib2.h>
#include <X11/Xlib.h>

#ifdef XINERAMA
#include <X11/extensions/Xinerama.h>
#endif

#define nil NULL

#define MIN(a, b) ((a) < (b) ? (a) : (b))
#define MAX(a, b) ((a) > (b) ? (a) : (b))
#define LENGTH(x) (sizeof x / sizeof x[0])

/* image modes */
typedef enum { ModeCenter, ModeZoom, ModeScale } ImageMode;

struct Monitor {
    int x, y, w, h;
};

static int sx, sy, sw, sh; /* screen geometry */
static Display* dpy;
static Window root;
static int nmonitor, nimage; /* Amount of monitors/available background
                                images */
static struct Monitor monitors[8];
static Imlib_Image images[LENGTH(monitors)];

/* free images before exit */
// returns nil or an error message
char* cleanup(void)
{
    for (int i = 0; i < nimage; i++) {
        imlib_context_set_image(images[i]);
        imlib_free_image_and_decache();
    }
    return nil;
}

/* draw background to root */
// returns nil or an error message
char* drawbg(ImageMode mode, bool rotate)
{
    int i, w, h, nx, ny, nh, nw, tmp;
    double factor;
    Pixmap pm;
    Imlib_Image tmpimg, buffer;

    pm = XCreatePixmap(dpy, root, sw, sh, DefaultDepth(dpy, DefaultScreen(dpy)));
    if (!(buffer = imlib_create_image(sw, sh))) {
        return "cannot allocate buffer";
    }
    imlib_context_set_image(buffer);
    imlib_image_fill_rectangle(0, 0, sw, sh);
    imlib_context_set_blend(1);
    for (i = 0; i < nmonitor; i++) {
        imlib_context_set_image(images[i % nimage]);
        w = imlib_image_get_width();
        h = imlib_image_get_height();
        if (!(tmpimg = imlib_clone_image())) {
            return "cannot clone image";
        }
        imlib_context_set_image(tmpimg);
        if (rotate
            && ((monitors[i].w > monitors[i].h && w < h)
                || (monitors[i].w < monitors[i].h && w > h))) {
            imlib_image_orientate(1);
            tmp = w;
            w = h;
            h = tmp;
        }
        imlib_context_set_image(buffer);
        switch (mode) {
        case ModeCenter:
            nw = (monitors[i].w - w) / 2;
            nh = (monitors[i].h - h) / 2;
            nx = monitors[i].x + (monitors[i].w - nw) / 2;
            ny = monitors[i].y + (monitors[i].h - nh) / 2;
            break;
        case ModeZoom:
            nw = monitors[i].w;
            nh = monitors[i].h;
            if (w > h && (w / h > (monitors[i].w / monitors[i].h))) {
                nx = monitors[i].x + (monitors[i].w - nw) / 2;
                ny = monitors[i].y + (int)ceil(h * nx / w) / 2;
            } else {
                ny = monitors[i].y + (monitors[i].h - nh) / 2;
                nx = monitors[i].x + (int)ceil(w * ny / h) / 2;
            }
            break;
        default: /* ModeScale */
            factor = MAX((double)w / monitors[i].w, (double)h / monitors[i].h);
            nw = w / factor;
            nh = h / factor;
            nx = monitors[i].x + (monitors[i].w - nw) / 2;
            ny = monitors[i].y + (monitors[i].h - nh) / 2;
        }
        imlib_blend_image_onto_image(tmpimg, 0, 0, 0, w, h, nx, ny, nw, nh);
        imlib_context_set_image(tmpimg);
        imlib_free_image();
    }
    imlib_context_set_blend(0);
    imlib_context_set_image(buffer);
    imlib_context_set_drawable(root);
    imlib_render_image_on_drawable(0, 0);
    imlib_context_set_drawable(pm);
    imlib_render_image_on_drawable(0, 0);
    XSetWindowBackgroundPixmap(dpy, root, pm);
    imlib_context_set_image(buffer);
    imlib_free_image_and_decache();
    XFreePixmap(dpy, pm);
}

/* update screen and/or Xinerama dimensions */
void updategeom(void)
{
#ifdef XINERAMA
    int i;
    XineramaScreenInfo* info = nil;

    if (XineramaIsActive(dpy) && (info = XineramaQueryScreens(dpy, &nmonitor))) {
        nmonitor = MIN(nmonitor, LENGTH(monitors));
        for (i = 0; i < nmonitor; i++) {
            monitors[i].x = info[i].x_org;
            monitors[i].y = info[i].y_org;
            monitors[i].w = info[i].width;
            monitors[i].h = info[i].height;
        }
        XFree(info);
    } else
#endif
    {
        nmonitor = 1;
        monitors[0].x = sx;
        monitors[0].y = sy;
        monitors[0].w = sw;
        monitors[0].h = sh;
    }
}

/* main loop */
// returns either nil or an error message
char* run(ImageMode mode, bool rotate)
{
    XEvent ev;

    // TODO: Remove the running feature and possibly also the X event include

    bool running = false;
    for (;;) {
        updategeom();
        drawbg(mode, rotate);
        if (!running) {
            break;
        }
        imlib_flush_loaders();
        XNextEvent(dpy, &ev);
        if (ev.type == ConfigureNotify) {
            sw = ev.xconfigure.width;
            sh = ev.xconfigure.height;
            imlib_flush_loaders();
        }
    }
    return nil;
}

// setup returns either nil or an error message
/* set up imlib and X */
char* setup(const char* filename, const char* col)
{
    Visual* vis;
    Colormap cm;
    XColor color;
    int i, screen;

    // TODO: Should multiple images be loaded for multiple monitors?
    /* Loading image */
    images[0] = imlib_load_image_without_cache(filename);
    nimage = 1;
    if (images[0] == nil) {
        return "no image to draw";
    }

    /* set up X */
    screen = DefaultScreen(dpy);
    vis = DefaultVisual(dpy, screen);
    cm = DefaultColormap(dpy, screen);
    root = RootWindow(dpy, screen);
    XSelectInput(dpy, root, StructureNotifyMask);
    sx = sy = 0;
    sw = DisplayWidth(dpy, screen);
    sh = DisplayHeight(dpy, screen);

    if (!XAllocNamedColor(dpy, cm, col, &color, &color)) {
        return "cannot allocate color";
    }

    /* set up Imlib */
    imlib_context_set_display(dpy);
    imlib_context_set_visual(vis);
    imlib_context_set_colormap(cm);
    imlib_context_set_color(color.red, color.green, color.blue, 255);

    return nil;
}

// SetBackground returns either nil or an error message
char* SetBackground(const char* filename, bool rotate, ImageMode mode, bool verbose)
{
    // TODO: Never die, return an error instead
    if (!(dpy = XOpenDisplay(nil))) {
        return "cannot open display";
    }

    const char* col = "#000000";
    char* err = setup(filename, col);
    if (err != nil) {
        return err;
    }
    err = run(mode, rotate);
    if (err != nil) {
        return err;
    }
    err = cleanup();
    if (err != nil) {
        return err;
    }
    XCloseDisplay(dpy);
    return nil;
}
