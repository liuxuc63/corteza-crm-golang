<template>
  <b-container
    fluid="xl"
    class="d-flex flex-column py-3"
  >
    <portal to="topbar-title">
      {{ $t('navigation.page') }}
    </portal>

    <b-card
      no-body
      class="shadow-sm h-100"
    >
      <b-card-header
        header-bg-variant="white"
        class="border-bottom"
      >
        <b-row
          no-gutters
          class="justify-content-between wrap-with-vertical-gutters"
        >
          <div class="flex-grow-1">
            <b-input-group
              v-if="namespace.canCreatePage"
              class="h-100"
            >
              <b-input
                id="name"
                v-model="page.title"
                data-test-id="input-name"
                required
                type="text"
                class="h-100"
                :placeholder="$t('newPlaceholder')"
              />
              <b-input-group-append>
                <b-button
                  data-test-id="button-create-page"
                  type="submit"
                  variant="primary"
                  size="lg"
                  @click="handleAddPageFormSubmit"
                >
                  {{ $t('createLabel') }}
                </b-button>
              </b-input-group-append>
            </b-input-group>
          </div>
          <div class="d-flex justify-content-sm-end flex-grow-1">
            <c-permissions-button
              v-if="namespace.canGrant"
              :resource="`corteza::compose:page/${namespace.namespaceID}/*`"
              class="btn-lg"
              :button-label="$t('label.permissions')"
              button-variant="light"
            />
          </div>
        </b-row>
        <b-row
          class="align-content-center mt-2"
        >
          <b-col
            cols="12"
          >
            <span class="text-muted font-italic">
              {{ $t('instructions') }}
            </span>
          </b-col>
        </b-row>
      </b-card-header>

      <div
        v-if="processing"
        class="text-center text-muted m-5"
      >
        <div>
          <b-spinner
            class="align-middle m-2"
          />
        </div>
        {{ $t('loading') }}
      </div>

      <page-tree
        v-else
        v-model="tree"
        :namespace="namespace"
        class="card overflow-auto h-100"
        @reorder="handleReorder"
      />
    </b-card>
  </b-container>
</template>

<script>
import { mapActions } from 'vuex'
import PageTree from 'corteza-webapp-compose/src/components/Admin/Page/Tree'
import { compose } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'page',
  },

  name: 'PageList',

  components: {
    PageTree,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: false,
      default: undefined,
    },
  },

  data () {
    return {
      tree: [],
      page: new compose.Page({ visible: true }),
      processing: false,
    }
  },

  created () {
    this.loadTree()
  },

  methods: {
    ...mapActions({
      createPage: 'page/create',
    }),

    loadTree () {
      this.processing = true
      const { namespaceID } = this.namespace
      this.$ComposeAPI.pageTree({ namespaceID }).then((tree) => {
        this.tree = tree.map(p => new compose.Page(p))
      }).catch(this.toastErrorHandler(this.$t('notification:page.loadFailed')))
        .finally(() => {
          this.processing = false
        })
    },

    handleAddPageFormSubmit () {
      const { namespaceID } = this.namespace
      this.page.weight = this.tree.length
      this.createPage({ ...this.page, namespaceID }).then(({ pageID }) => {
        this.$router.push({ name: 'admin.pages.edit', params: { pageID } })
      }).catch(this.toastErrorHandler(this.$t('notification:page.saveFailed')))
    },

    handleReorder () {
      this.loadTree()
    },
  },
}
</script>
<style lang="scss">
//!important usage to over-ride library styling
$input-height: 42px;
$content-height: 48px;
$blank-li-height: 10px;
$left-padding: 5px;
$border-color: $light;
$hover-color: $gray-200;
$dropping-color: $secondary;

.page-name-input {
  height: $input-height;
}

.sortable-tree {
  .content {
    height: 0px !important;
  }

  ul {
    .content {
      height: 100% !important;
      min-height: $content-height !important;
      line-height: $content-height !important;

      &:hover {
        background: $hover-color;
      }
    }
  }

  li {
    white-space: nowrap;
    background: white;

    &.blank-li {
      height: $blank-li-height !important;

      .sortable-tree {
        max-height: 100%;
      }

      &:nth-last-of-type(1)::before {
        border-left-color: white !important;
        height: 0;
      }
    }

    &::before {
      top: $content-height / -2 !important;
      border-left-color: white !important;
    }

    &::after {
      height: $content-height !important;
      top: $content-height / 2 !important;
      border-color: #fff !important;
    }

    &.parent-li:nth-last-child(2)::before {
      height: $content-height !important;
      top: $content-height / -2 !important;
    }
  }

  .parent-li {
    border-top: 1px solid $border-color;

    .exist-li, .blank-li {
      border-top: none;

      &::after {
        border-top: 2px solid $border-color !important;
        margin-left: 0;
      }

      &::before {
        border-left: 2px solid $border-color !important;
      }
    }

    &.blank-li {
      &::before {
        border-left: 2px solid $border-color !important;
      }
    }

    &.exist-li {
      &::before {
        border-color: white !important;
      }

      .parent-li {
        &.exist-li {
          &::before {
            border-color: $border-color !important;
          }
        }
      }
    }
  }
}

.droper {
  background: $dropping-color !important;
}

.pages-list-header {
  min-height: $content-height;
  background-color: $gray-200;
  margin-bottom: -1.8rem !important;
  border-bottom: 2px solid $light;
  z-index: 1;
}

</style>
