# These are some common variables for Make

IMG_TAG ?= latest

# Image URL to use all building/pushing image targets
QUAY_IMG ?= quay.io/ccremer/greposync:$(IMG_TAG)