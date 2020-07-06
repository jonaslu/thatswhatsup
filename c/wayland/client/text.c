#include <cairo/cairo.h>

void render_text(unsigned char *buffer, int width, int height, int stride)
{
  cairo_surface_t *cairo_surface = cairo_image_surface_create_for_data(buffer, CAIRO_FORMAT_ARGB32, width, height, width * stride);
  cairo_t *cairo = cairo_create(cairo_surface);

  cairo_select_font_face(cairo, "Roboto Mono", CAIRO_FONT_SLANT_NORMAL, CAIRO_FONT_WEIGHT_BOLD);
  cairo_set_font_size(cairo, 32.0);

  cairo_move_to(cairo, 100.0, 50.0);
  cairo_set_source_rgba(cairo, 1.0, 1.0, 1.0, 1.0);

  cairo_show_text(cairo, "Cheese and goats");

  cairo_destroy(cairo);
  cairo_surface_destroy(cairo_surface);
}