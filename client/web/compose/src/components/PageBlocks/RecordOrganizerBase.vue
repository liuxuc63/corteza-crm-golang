<template>
  <wrap
    v-bind="$props"
    v-on="$listeners"
  >
    <template #default>
      <div
        v-if="!isConfigured"
        class="d-flex align-items-center justify-content-center h-100 p-2"
      >
        {{ $t('recordOrganizer.notConfigured') }}
      </div>
      <div
        v-else
        class="h-100"
      >
        <draggable
          v-if="!processing"
          :id="draggableID"
          v-model="records"
          handle=".record-item"
          :group="{ name: moduleID, pull: canPull, put: canPut }"
          :move="checkMove"
          class="h-100 pt-2 px-2 overflow-auto"
          @change="onChange"
        >
          <template
            v-if="!records.length"
            #header
          >
            <div
              class="small p-2 text-secondary"
            >
              {{ $t('recordOrganizer.noRecords') }}
            </div>
          </template>
          <router-link
            v-for="record in records"
            :key="`${record.recordID}`"
            tag="div"
            class="record-item mb-2"
            :class="{ 'pointer': roRecordPage }"
            :to="{ name: 'page.record', params: { pageID: (roRecordPage || {}).pageID, recordID: record.recordID }, query: null }"
          >
            <b-card
              body-class="px-2 py-1"
              class="border rounded"
            >
              <b-card-title
                v-if="labelField"
                title-tag="div"
                class="mb-1"
              >
                <field-viewer
                  v-if="labelField.canReadRecordValue"
                  :field="labelField"
                  :record="record"
                  :namespace="namespace"
                  value-only
                />
                <i
                  v-else
                  class="text-secondary"
                >{{ $t('field.noPermission') }}</i>
              </b-card-title>
              <b-card-text
                v-if="descriptionField"
                class="small"
              >
                <field-viewer
                  v-if="descriptionField.canReadRecordValue"
                  :field="descriptionField"
                  :record="record"
                  :namespace="namespace"
                  value-only
                />
                <i
                  v-else
                  class="text-secondary"
                >
                  {{ $t('field.noPermission') }}
                </i>
              </b-card-text>
            </b-card>
          </router-link>
        </draggable>
        <div
          v-else
          class="d-flex align-items-center justify-content-center h-100"
        >
          <b-spinner />
        </div>
      </div>
    </template>
    <template
      v-if="canAddRecord"
      #footer
    >
      <div
        class="p-2"
      >
        <b-button
          variant="outline-light"
          class="d-flex align-items-center justify-content-center text-primary border-0"
          @click.prevent="createNewRecord"
        >
          {{ $t('recordOrganizer.addNewRecord') }}
        </b-button>
      </div>
    </template>
  </wrap>
</template>

<script>
import { mapGetters } from 'vuex'
import base from './base'
import draggable from 'vuedraggable'
import FieldViewer from 'corteza-webapp-compose/src/components/ModuleFields/Viewer'
import users from 'corteza-webapp-compose/src/mixins/users'
import { evaluatePrefilter, getFieldFilter } from 'corteza-webapp-compose/src/lib/record-filter'
import { compose, NoID } from '@cortezaproject/corteza-js'

export default {
  i18nOptions: {
    namespaces: 'block',
  },

  components: {
    FieldViewer,
    draggable,
  },

  extends: base,

  mixins: [
    users,
  ],

  data () {
    return {
      processing: false,

      filter: {
        sort: '',
        query: '',
      },

      records: [],
    }
  },

  computed: {
    ...mapGetters({
      getModuleByID: 'module/getByID',
      pages: 'page/set',
    }),

    draggableID () {
      return `recordOrganizer-${this.blockIndex}`
    },

    roModule () {
      return this.getModuleByID(this.moduleID)
    },

    roRecordPage () {
      return this.pages.find(p => p.moduleID === this.moduleID)
    },

    moduleID () {
      return this.options.moduleID
    },

    labelField () {
      const { labelField } = this.options

      if (!labelField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === labelField) || {}
    },

    descriptionField () {
      const { descriptionField } = this.options

      if (!descriptionField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === descriptionField) || {}
    },

    positionField () {
      const { positionField } = this.options

      if (!positionField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === positionField) || {}
    },

    groupField () {
      const { groupField } = this.options

      if (!groupField) {
        return undefined
      }

      return this.roModule.fields.find(f => f.name === groupField) || {}
    },

    canPull () {
      return this.positionField ? this.positionField.canUpdateRecordValue : true
    },

    canPut () {
      return this.canPull && (this.groupField ? this.groupField.canUpdateRecordValue : true)
    },

    canAddRecord () {
      return this.roModule && this.roModule.canCreateRecord && this.roRecordPage
    },

    isConfigured () {
      return !!(this.labelField || this.descriptionField)
    },
  },

  watch: {
    options: {
      deep: true,
      handler () {
        this.refresh()
      },
    },

    'record.recordID': {
      immediate: true,
      handler () {
        this.refresh()
      },
    },
  },

  mounted () {
    this.$root.$on(`refetch-non-record-blocks:${this.page.pageID}`, this.refresh)
  },

  beforeDestroy () {
    this.$root.$off(`refetch-non-record-blocks:${this.page.pageID}`)
  },

  methods: {
    // Allow move if repositioned or if record isn't in target record organizer
    checkMove ({ draggedContext = {}, relatedContext = {} }) {
      const { moduleID, recordID } = draggedContext.element || {}
      const { $attrs = {}, $el = {}, $options = {} } = relatedContext.component || {}
      const relatedRecords = ($options.propsData || {}).value || []

      if (moduleID !== $attrs.group.name) {
        return false
      }

      return this.draggableID === $el.id || !relatedRecords.some(r => r.recordID === recordID)
    },

    onChange ({ added, moved }) {
      if (added) {
        this.reorganize(added)
      } else if (moved) {
        this.reposition(moved)
      }
    },

    reorganize ({ element: record, newIndex }) {
      // Move record to a different position in a different group
      this.moveRecord(
        record,
        this.calcNewPosition(record, newIndex),
        this.options.group,
      )
    },

    reposition ({ element: record, newIndex }) {
      // Move record to a different position in the same group
      this.moveRecord(
        record,
        this.calcNewPosition(record, newIndex),
      )
    },

    /**
     * Calculates optimal position value for dropped record
     */
    calcNewPosition (record, newPosition = 0) {
      if (newPosition <= 0) {
        // Dropped in first place, easy-breezy
        return 0
      }

      const total = this.records.length
      if (newPosition > total) {
        // Dropped at the end,
        // make sure we don't put it too far away
        return total
      }

      // Find position field on the record placed before the drop position
      // fallback to 1
      return parseInt(this.records[newPosition - 1].values[this.options.positionField] || 0) + 1
    },

    createNewRecord () {
      const { groupField, group } = this.options

      if (!this.roRecordPage) {
        // can not create record without a record page
        return
      }

      const { pageID } = this.roRecordPage

      // Prefill values with the group value set in the options
      const values = {}
      if (groupField) {
        values[groupField] = group
      }

      this.$router.push({ name: 'page.record.create', params: { pageID, values: values, refRecord: this.record } })
    },

    expandFilter () {
      const filter = []

      /* eslint-disable no-template-curly-in-string */
      if (!this.record) {
        // If there is no current record and we are using recordID/ownerID variable in (pre)filter
        // we should disable the block
        if ((this.options.filter || '').includes('${record')) {
          throw Error(this.$t('notification:record.invalidRecordVar'))
        }

        if ((this.options.filter || '').includes('${ownerID}')) {
          throw Error(this.$t('notification:record.invalidOwnerVar'))
        }
      }

      if (this.options.filter) {
        // Little magic here: filter is wraped with backticks and evaluated
        // this allows us to us ${record.values....}, ${recordID}, ${ownerID}, ${userID} in filter string;
        // hence the /hanging/ record, recordID, ownerID and userID variables
        filter.push(`(${evaluatePrefilter(this.options.filter, {
          record: this.record,
          recordID: (this.record || {}).recordID || NoID,
          ownerID: (this.record || {}).ownedBy || NoID,
          userID: (this.$auth.user || {}).userID || NoID,
        })})`)
      }

      if (this.groupField && this.options.group !== undefined) {
        filter.push(`(${getFieldFilter(this.groupField.name, this.groupField.kind, this.options.group, '=')})`)
      }

      return filter.join(' AND ')
    },

    /**
     * Reposition and optionally move record to a different group
     *
     * This is only a helper function and we do not keep any hard dependencies on
     * the API client.
     *
     * @param {Compose}           api Compose API client
     * @param {Record}            record,     Record we're moving
     * @param {Number}            position    New position
     * @param {String|undefined}  group       New group
     * @returns {Promise<void>}
     */
    async moveRecord (record, position, group) {
      const { namespaceID, moduleID, recordID } = record

      if (moduleID !== this.options.moduleID) {
        throw Error(this.$t('record.moduleMismatch'))
      }

      const { positionField, groupField } = this.options
      const args = {
        recordID,
        filter: this.expandFilter(),
        positionField,
        position,
      }

      if (group !== undefined) {
        // If group is set (empty string is a valid!
        args.groupField = groupField
        args.group = group || ''
      }

      const params = {
        procedure: 'organize',
        namespaceID,
        moduleID,
        // map kv to [{ name: k, value: v }, ...]
        args: Object.keys(args).map(name => ({ name, value: String(args[name]) })),
      }

      return this.$ComposeAPI.recordExec(params).then(this.fetchRecords)
    },

    /**
     * Fetches group of records using configured options & module
     *
     * @param {Module}            module Module to use for assembling API request & casting results
     * @param {String}            query Filter records
     * @returns {Promise<Record[]>}
     */
    async fetchRecords () {
      if (!this.roModule) {
        return
      }

      if (this.roModule.moduleID !== this.options.moduleID) {
        throw Error(this.$t('record.moduleMismatch'))
      }

      const query = this.expandFilter()
      const { labelField, descriptionField, positionField } = this.options
      const { moduleID, namespaceID } = this.roModule
      const sort = positionField || `updatedAt, ${labelField || descriptionField}`

      this.processing = true

      return this.$ComposeAPI.recordList({ namespaceID, moduleID, query, sort })
        .then(({ set }) => {
          this.records = set.map(r => Object.freeze(new compose.Record(this.roModule, r)))
          const fields = [this.labelField, this.descriptionField].filter(f => !!f)
          this.fetchUsers(fields, this.records)
        }).catch(e => {
          console.error(e)
        }).finally(() => {
          this.processing = false
        })
    },

    refresh () {
      this.fetchRecords()
    },
  },
}
</script>

<style lang="scss" scoped>
.grab {
  cursor: grab !important;
}

.record-item:hover {
  box-shadow: 0 0.125rem 0.25rem rgb(0 0 0 / 8%) !important;
}
</style>
