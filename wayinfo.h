#pragma once

#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <wayland-client.h>

/* This file uses a buffer that is estimated to be large enough to collect all needed information */

#define BUF buf + strlen(buf)
#define RETBUF                                                                                    \
    char* str = calloc(strlen(buf) + 1, sizeof(char));                                            \
    strcpy(str, buf);                                                                             \
    return str

static char buf[10240];

typedef void (*sprint_info_t)(void* info);

struct global_info {
    struct wl_list link;

    uint32_t id;
    uint32_t version;
    char* interface;

    sprint_info_t sprint;
};

struct output_mode {
    struct wl_list link;

    uint32_t flags;
    int32_t width, height;
    int32_t refresh;
};

struct output_info {
    struct global_info global;

    struct wl_output* output;

    struct {
        int32_t x, y;
        int32_t physical_width, physical_height;
        enum wl_output_subpixel subpixel;
        enum wl_output_transform output_transform;
        char* make;
        char* model;
    } geometry;

    struct wl_list modes;
};

struct shm_format {
    struct wl_list link;

    uint32_t format;
};

struct shm_info {
    struct global_info global;
    struct wl_shm* shm;

    struct wl_list formats;
};

struct seat_info {
    struct global_info global;
    struct wl_seat* seat;

    uint32_t capabilities;
};

struct weston_info {
    struct wl_display* display;
    struct wl_registry* registry;

    struct wl_list infos;
    bool roundtrip_needed;
};

void sprint_global_info(void* data)
{
    struct global_info* global = data;

    sprintf(BUF, "interface: '%s', version: %u, name: %u\n", global->interface, global->version,
        global->id);
}

static void init_global_info(struct weston_info* info, struct global_info* global, uint32_t id,
    const char* interface, uint32_t version)
{
    global->id = id;
    global->version = version;
    global->interface = strdup(interface);

    wl_list_insert(info->infos.prev, &global->link);
}

void sprint_output_info(void* data)
{
    struct output_info* output = data;
    struct output_mode* mode;
    const char* subpixel_orientation;
    const char* transform;

    sprint_global_info(data);

    switch (output->geometry.subpixel) {
    case WL_OUTPUT_SUBPIXEL_UNKNOWN:
        subpixel_orientation = "unknown";
        break;
    case WL_OUTPUT_SUBPIXEL_NONE:
        subpixel_orientation = "none";
        break;
    case WL_OUTPUT_SUBPIXEL_HORIZONTAL_RGB:
        subpixel_orientation = "horizontal rgb";
        break;
    case WL_OUTPUT_SUBPIXEL_HORIZONTAL_BGR:
        subpixel_orientation = "horizontal bgr";
        break;
    case WL_OUTPUT_SUBPIXEL_VERTICAL_RGB:
        subpixel_orientation = "vertical rgb";
        break;
    case WL_OUTPUT_SUBPIXEL_VERTICAL_BGR:
        subpixel_orientation = "vertical bgr";
        break;
    default:
        fprintf(stderr, "unknown subpixel orientation %u\n", output->geometry.subpixel);
        subpixel_orientation = "unexpected value";
        break;
    }

    switch (output->geometry.output_transform) {
    case WL_OUTPUT_TRANSFORM_NORMAL:
        transform = "normal";
        break;
    case WL_OUTPUT_TRANSFORM_90:
        transform = "90°";
        break;
    case WL_OUTPUT_TRANSFORM_180:
        transform = "180°";
        break;
    case WL_OUTPUT_TRANSFORM_270:
        transform = "270°";
        break;
    case WL_OUTPUT_TRANSFORM_FLIPPED:
        transform = "flipped";
        break;
    case WL_OUTPUT_TRANSFORM_FLIPPED_90:
        transform = "flipped 90°";
        break;
    case WL_OUTPUT_TRANSFORM_FLIPPED_180:
        transform = "flipped 180°";
        break;
    case WL_OUTPUT_TRANSFORM_FLIPPED_270:
        transform = "flipped 270°";
        break;
    default:
        fprintf(stderr, "unknown output transform %u\n", output->geometry.output_transform);
        transform = "unexpected value";
        break;
    }

    sprintf(BUF, "\tx: %d, y: %d,\n", output->geometry.x, output->geometry.y);
    sprintf(BUF, "\tphysical_width: %d mm, physical_height: %d mm,\n",
        output->geometry.physical_width, output->geometry.physical_height);
    sprintf(BUF, "\tmake: '%s', model: '%s',\n", output->geometry.make, output->geometry.model);
    sprintf(BUF, "\tsubpixel_orientation: %s, output_tranform: %s,\n", subpixel_orientation,
        transform);

    wl_list_for_each(mode, &output->modes, link)
    {
        sprintf(BUF, "\tmode:\n");
        sprintf(BUF, "\t\twidth: %d px, height: %d px, refresh: %.f Hz,\n", mode->width,
            mode->height, (float)mode->refresh / 1000);
        sprintf(BUF, "\t\tflags:");
        if (mode->flags & WL_OUTPUT_MODE_CURRENT) {
            sprintf(BUF, " current");
        }
        if (mode->flags & WL_OUTPUT_MODE_PREFERRED) {
            sprintf(BUF, " preferred");
        }
        sprintf(BUF, "\n");
    }
}

void sprint_shm_info(void* data)
{
    struct shm_info* shm = data;
    struct shm_format* format;

    sprint_global_info(data);
    sprintf(BUF, "%s", "\tformats:");

    wl_list_for_each(format, &shm->formats, link)
        sprintf(BUF, " %s", (format->format == WL_SHM_FORMAT_ARGB8888) ? "ARGB8888" : "XRGB8888");

    sprintf(BUF, "\n");
}

void sprint_seat_info(void* data)
{
    struct seat_info* seat = data;

    sprint_global_info(data);
    sprintf(BUF, "%s", "\tcapabilities:");

    if (seat->capabilities & WL_SEAT_CAPABILITY_POINTER) {
        sprintf(BUF, " pointer");
    }
    if (seat->capabilities & WL_SEAT_CAPABILITY_KEYBOARD) {
        sprintf(BUF, " keyboard");
    }
    if (seat->capabilities & WL_SEAT_CAPABILITY_TOUCH) {
        sprintf(BUF, " touch");
    }

    sprintf(BUF, "\n");
}

static void seat_handle_capabilities(
    void* data, struct wl_seat* wl_seat, enum wl_seat_capability caps)
{
    struct seat_info* seat = data;
    seat->capabilities = caps;
}

static const struct wl_seat_listener seat_listener = {
    seat_handle_capabilities,
};

static void add_seat_info(struct weston_info* info, uint32_t id, uint32_t version)
{
    struct seat_info* seat = malloc(sizeof *seat);

    init_global_info(info, &seat->global, id, "wl_seat", version);
    seat->global.sprint = sprint_seat_info;

    seat->seat = wl_registry_bind(info->registry, id, &wl_seat_interface, 1);
    wl_seat_add_listener(seat->seat, &seat_listener, seat);

    info->roundtrip_needed = true;
}

static void shm_handle_format(void* data, struct wl_shm* wl_shm, uint32_t format)
{
    struct shm_info* shm = data;
    struct shm_format* shm_format = malloc(sizeof *shm_format);

    wl_list_insert(&shm->formats, &shm_format->link);
    shm_format->format = format;
}

static const struct wl_shm_listener shm_listener = {
    shm_handle_format,
};

static void add_shm_info(struct weston_info* info, uint32_t id, uint32_t version)
{
    struct shm_info* shm = malloc(sizeof *shm);

    init_global_info(info, &shm->global, id, "wl_shm", version);
    shm->global.sprint = sprint_shm_info;
    wl_list_init(&shm->formats);

    shm->shm = wl_registry_bind(info->registry, id, &wl_shm_interface, 1);
    wl_shm_add_listener(shm->shm, &shm_listener, shm);

    info->roundtrip_needed = true;
}

static void output_handle_geometry(void* data, struct wl_output* wl_output, int32_t x, int32_t y,
    int32_t physical_width, int32_t physical_height, int32_t subpixel, const char* make,
    const char* model, int32_t output_transform)
{
    struct output_info* output = data;

    output->geometry.x = x;
    output->geometry.y = y;
    output->geometry.physical_width = physical_width;
    output->geometry.physical_height = physical_height;
    output->geometry.subpixel = subpixel;
    output->geometry.make = strdup(make);
    output->geometry.model = strdup(model);
    output->geometry.output_transform = output_transform;
}

static void output_handle_mode(void* data, struct wl_output* wl_output, uint32_t flags,
    int32_t width, int32_t height, int32_t refresh)
{
    struct output_info* output = data;
    struct output_mode* mode = malloc(sizeof *mode);

    mode->flags = flags;
    mode->width = width;
    mode->height = height;
    mode->refresh = refresh;

    wl_list_insert(output->modes.prev, &mode->link);
}

static const struct wl_output_listener output_listener = {
    output_handle_geometry,
    output_handle_mode,
};

static void add_output_info(struct weston_info* info, uint32_t id, uint32_t version)
{
    struct output_info* output = malloc(sizeof *output);

    init_global_info(info, &output->global, id, "wl_output", version);
    output->global.sprint = sprint_output_info;

    wl_list_init(&output->modes);

    output->output = wl_registry_bind(info->registry, id, &wl_output_interface, 1);
    wl_output_add_listener(output->output, &output_listener, output);

    info->roundtrip_needed = true;
}

static void add_global_info(
    struct weston_info* info, uint32_t id, const char* interface, uint32_t version)
{
    struct global_info* global = malloc(sizeof *global);

    init_global_info(info, global, id, interface, version);
    global->sprint = sprint_global_info;
}

static void global_handler(
    void* data, struct wl_registry* registry, uint32_t id, const char* interface, uint32_t version)
{
    struct weston_info* info = data;

    if (!strcmp(interface, "wl_seat"))
        add_seat_info(info, id, version);
    else if (!strcmp(interface, "wl_shm"))
        add_shm_info(info, id, version);
    else if (!strcmp(interface, "wl_output"))
        add_output_info(info, id, version);
    else
        add_global_info(info, id, interface, version);
}

static void global_remove_handler(void* data, struct wl_registry* registry, uint32_t name) {}

static const struct wl_registry_listener registry_listener
    = { global_handler, global_remove_handler };

void sprint_infos(struct wl_list* infos)
{
    struct global_info* info;
    wl_list_for_each(info, infos, link) info->sprint(info);
}

bool WaylandRunning()
{
    struct wl_display* dpy = wl_display_connect(NULL);
    bool canConnect = (bool)dpy;
    if (canConnect) {
        wl_display_disconnect(dpy);
    }
    return canConnect;
}

char* WaylandInfoString()
{
    buf[0] = 0;
    buf[1] = 0;

    struct weston_info info;

    info.display = wl_display_connect(NULL);
    if (!info.display) {
        sprintf(BUF, "wayland plugin info string: wayland is not in use right now\n");
        RETBUF;
    }

    wl_list_init(&info.infos);

    info.registry = wl_display_get_registry(info.display);
    wl_registry_add_listener(info.registry, &registry_listener, &info);

    do {
        info.roundtrip_needed = false;
        wl_display_roundtrip(info.display);
    } while (info.roundtrip_needed);

    sprint_infos(&info.infos);

    if (info.display) {
        wl_display_disconnect(info.display);
    }

    char* str = calloc(strlen(buf) + 1, sizeof(char));
    strcpy(str, buf);
    buf[0] = 0;
    buf[1] = 0;
    return str;
}
