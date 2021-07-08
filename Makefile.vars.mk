# These are some common variables for Make

IMG_TAG ?= latest

# Image URL to use all building/pushing image targets
QUAY_IMG ?= quay.io/ccremer/greposync:$(IMG_TAG)

GOASCIIDOC_OUT_PATH ?= docs/modules/ROOT/pages/references/godoc.adoc

GOASCIIDOC_ARGS_BUILD ?= -o $(GOASCIIDOC_OUT_PATH) --templatedir docs/godoc-templates
GOASCIIDOC_CMD = go run github.com/mariotoffia/goasciidoc $(GOASCIIDOC_ARGS_BUILD)
