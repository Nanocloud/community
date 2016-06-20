import Ember from 'ember';

export default Ember.Component.extend({

  mode: null,
  placeholder: "",
  errorMessage: "",
  value: "",
  focus: false,
  isValid: true,
  type: "text",

  didInsertElement() {
    if (this.get('autofocus') === true) {
      this.$().find('input').focus();
    }
  },

  init: function() {
    this._super(...arguments);
    var valuePath = this.get('valuePath');

    if (this.get('model')) {
      Ember.defineProperty(this, 'value', Ember.computed.alias(`model.${valuePath}`));
    }
  },

  fieldIsCorrect: Ember.computed('isValid', 'value', function() {
    if (this.get('isValid') === true && this.get('value') !== undefined && this.get('value') !== "") {
      return true;
    }
    return false;
  }),

  getErrorMessage: function() {

    if (this.get('model')) {
      if (this.get('value')) {
        setTimeout(function() {
          var errorMessage = this.get('model').get('validations.attrs').get(this.get('valuePath')).get('message');
          this.set('errorMessage', errorMessage);
          this.set('isValid', Ember.isEmpty(errorMessage));
        }.bind(this), 500);
      }
    }
  }.observes('value', 'focus'),

  actions: {
    toggleFocus: function() {
      this.toggleProperty('focus');
    }
  }
});
