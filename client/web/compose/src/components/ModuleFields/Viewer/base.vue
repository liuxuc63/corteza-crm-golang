<template>
  <div>
    <!-- Extra empty line is added thanks to white-space: pre-line (multivalue) if we write div in multiple lines  -->
    <!-- eslint-disable-next-line -->
    <div :class="classes">{{ formatted }}</div>

    <errors :errors="errors" />
  </div>
</template>
<script>
import errors from '../errors'
import { compose, validator } from '@cortezaproject/corteza-js'

export default {
  components: {
    // errors is used in the components that extends base
    // eslint-disable-next-line vue/no-unused-components
    errors,
  },

  props: {
    namespace: {
      type: compose.Namespace,
      required: true,
    },

    field: {
      type: compose.ModuleField,
      required: true,
    },

    record: {
      type: compose.Record,
      required: true,
    },

    errors: {
      type: validator.Validated,
      default () { return new validator.Validated() },
    },

    valueOnly: {
      type: Boolean,
      required: false,
    },
  },

  computed: {
    value () {
      if (this.field.isSystem) {
        return this.record[this.field.name]
      }

      return this.record ? this.record.values[this.field.name] : undefined
    },

    formatted () {
      if (this.field.isMulti) {
        return this.value.join(this.field.options.multiDelimiter)
      }
      return this.value
    },

    classes () {
      if (this.field.isMulti) {
        return ['multiline']
      }
      return []
    },
  },
}
</script>

<style>
.multiline {
  white-space: pre-line;
}
</style>
