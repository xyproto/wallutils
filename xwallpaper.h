#include "X11/bitmaps/gray"
#include <X11/Xatom.h>
#include <X11/Xcursor/Xcursor.h>
#include <X11/Xlib.h>
#include <X11/Xmu/CurUtil.h>
#include <X11/Xutil.h>
#include <X11/xpm.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

static Display* dpy;
static int screen;
static Window root;

// Read black/white X bitmap file
static Pixmap ReadXBMFile(char* filename, unsigned int* width, unsigned int* height, int* x_hot, int* y_hot, bool verbose)
{
    Pixmap bitmap;
    int status = XReadBitmapFile(dpy, root, filename, width, height, &bitmap, x_hot, y_hot);
    if (status == BitmapSuccess) {
        return (bitmap);
    }
    if (status == BitmapOpenFailed) {
        fprintf(stderr, "can't open file: %s\n", filename);
    } else if (status == BitmapFileInvalid) {
        fprintf(stderr, "bad bitmap format file: %s\n", filename);
    }
    fprintf(stderr, "insufficient memory for bitmap: %s", filename);
    exit(1);
}

// Read colorful X pixmap file
static Pixmap ReadXPMFile(char* filename, unsigned int* width, unsigned int* height, bool verbose)
{
    Pixmap pixmap;
    XpmAttributes attributes;
    // shapemask is set to NULL and is only used for transparency
    int status = XpmReadFileToPixmap(dpy, root, filename, &pixmap, NULL, &attributes);
    if (status == XpmSuccess) {
        *width = attributes.width;
        *height = attributes.height;
        return (pixmap);
    }
    fprintf(stderr, "could not read XPM: %s", filename);
    exit(1);
}

// Set the background image to the given pixmap
static void SetBackgroundToPixmapAndFree(Pixmap bitmap, unsigned int width, unsigned int height, bool verbose)
{
    XGCValues gc_init;

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XCreateGC\n");
    GC gc = XCreateGC(dpy, root, GCForeground | GCBackground, &gc_init);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XCreatePixmap\n");
    Pixmap pix = XCreatePixmap(dpy, root, width, height, (unsigned int)DefaultDepth(dpy, screen));

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XCopyPlane\n");
    XCopyPlane(dpy, bitmap, pix, gc, 0, 0, width, height, 0, 0, (unsigned long)1);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XSetWindowBackgroundPixmap\n");
    XSetWindowBackgroundPixmap(dpy, root, pix);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XFreeGC\n");
    XFreeGC(dpy, gc);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XFreePixmap(dpy, bitmap)\n");
    XFreePixmap(dpy, bitmap);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XFreePixmap(dpy, pix)\n");
    XFreePixmap(dpy, pix);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree: XClearWindow\n");
    XClearWindow(dpy, root);
}

int SetBackground(char* filename, bool verbose)
{
    unsigned int ww;
    unsigned int hh;

    if (verbose) printf("X11: SetBackground, filename = %s\n", filename);

    if (verbose) printf("X11: XOpenDisplay... ");
    dpy = XOpenDisplay("");
    if (!dpy) {
        if (verbose) printf("failed\n");
        return -1;
    }
    if (verbose) printf("ok\n");

    screen = DefaultScreen(dpy);
    root = RootWindow(dpy, screen);

    if (verbose) printf("X11: root window OK\n");

    //Pixmap bitmap = ReadBitmapFile(filename, &ww, &hh, (int*)NULL, (int*)NULL, verbose);

    if (verbose) printf("X11: ReadXPMFile\n");
    Pixmap bitmap = ReadXPMFile(filename, &ww, &hh, verbose);

    if (verbose) printf("X11: SetBackgroundToPixmapAndFree\n");
    SetBackgroundToPixmapAndFree(bitmap, ww, hh, verbose);

    if (dpy) {
        if (verbose)  printf("X11: XCloseDisplay\n");
        XCloseDisplay(dpy);
    }

    return 0;
}
