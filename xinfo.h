#ifdef WIN32
#include <X11/Xwindows.h>
#endif
#include <X11/Xlib.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* This file uses a buffer that is estimated to be large enough to collect all
 * required information.
 *
 * Also, thanks to the author of this answer: https://stackoverflow.com/a/48843306/131264
 */

#define BUF buf + strlen(buf)
#define RETBUF                                                                                    \
    char* str = calloc(strlen(buf) + 1, sizeof(char));                                            \
    strcpy(str, buf);                                                                             \
    return str

static char buf[10240];

static void sprint_screen_info(Display* dpy, int scr)
{
    /*
     * there are 2.54 centimeters to an inch; so there are 25.4 millimeters.
     *
     *     dpi = N pixels / (M millimeters / (25.4 millimeters / 1 inch))
     *         = N pixels / (M inch / 25.4)
     *         = N * 25.4 pixels / M inch
     */

    double xres, yres;

    xres = ((((double)DisplayWidth(dpy, scr)) * 25.4) / ((double)DisplayWidthMM(dpy, scr)));
    yres = ((((double)DisplayHeight(dpy, scr)) * 25.4) / ((double)DisplayHeightMM(dpy, scr)));

    sprintf(BUF, "\n");
    sprintf(BUF, "screen #%d:\n", scr);
    sprintf(BUF, "  dimensions:    %dx%d pixels (%dx%d millimeters)\n", XDisplayWidth(dpy, scr),
        XDisplayHeight(dpy, scr), XDisplayWidthMM(dpy, scr), XDisplayHeightMM(dpy, scr));
    sprintf(BUF, "  resolution:    %dx%d dots per inch\n", (int)(xres + 0.5), (int)(yres + 0.5));
}

bool X11Running()
{
    Display* dpy = XOpenDisplay(NULL);
    bool canConnect = (bool)dpy;
    if (canConnect) {
        XCloseDisplay(dpy);
    }
    return canConnect;
}

char* X11InfoString()
{
    buf[0] = 0;
    buf[1] = 0;

    Display* dpy; /* X connection */
    char* displayname = NULL; /* server to contact */
    int i;

    dpy = XOpenDisplay(displayname);
    if (!dpy) {
        sprintf(BUF, "unable to open display \"%s\".\n", XDisplayName(displayname));
        RETBUF;
    }

    sprintf(BUF, "name of display:    %s\n", DisplayString(dpy));
    sprintf(BUF, "default screen number:    %d\n", DefaultScreen(dpy));
    sprintf(BUF, "number of screens:    %d\n", ScreenCount(dpy));

    for (i = 0; i < ScreenCount(dpy); i++) {
        sprintf(BUF, "SCREEN %d\n", i);
        sprint_screen_info(dpy, i);
    }

    if (dpy) {
        XCloseDisplay(dpy);
    }
    char* str = calloc(strlen(buf) + 1, sizeof(char));
    strcpy(str, buf);
    buf[0] = 0;
    buf[1] = 0;
    return str;
}
