all: doc
clean: clean-doc

EXAMPLE_DIR = docs/examples
EXAMPLE_PACKAGES = \
	async-promise \
	callback \
	describe-args \
	hello-world \
	js

NAPI_LIB_SUFFIX = .node

TARGET_BUILDDIR = build

EXAMPLE_BINDINGS = $(addsuffix $(NAPI_LIB_SUFFIX),$(EXAMPLE_PACKAGES))
TARGET_EXAMPLES = \
  $(addprefix $(TARGET_BUILDDIR)/, $(EXAMPLE_BINDINGS))

# TODO: Configure CGO_LDFLAGS_ALLOW for non-darwin systems.
CGO_LDFLAGS_ALLOW = (-Wl,(-undefined,dynamic_lookup|-no_pie|-search_paths_first))

doc: $(TARGET_EXAMPLES)

$(TARGET_EXAMPLES): | $(TARGET_BUILDDIR)
$(TARGET_EXAMPLES): $(TARGET_BUILDDIR)/%$(NAPI_LIB_SUFFIX): $(EXAMPLE_DIR)/%
	CGO_LDFLAGS_ALLOW='$(CGO_LDFLAGS_ALLOW)' \
	  go build -buildmode=c-shared -o "$(@)" "./$(<)/"

$(TARGET_BUILDDIR):
	mkdir -p "$(TARGET_BUILDDIR)"

clean:
	rmdir "$(TARGET_BUILDDIR)"

clean-doc:
	rm -f $(patsubst %,"%",$(TARGET_EXAMPLES))

.PHONY: all doc
.PHONY: clean clean-doc
