import Ember from 'ember';

export default Ember.Component.extend({

    isEditing: false,

    getInputType: function() {
      if (this.get('hideInput')) {
        return "password";
      }
      else {
        return "text";
      }
    }.property('hideInput'),

    actions: {

      toggle() {
        console.log('toggle');
        this.toggleProperty('isEditing');
      },

      submit() {
        this.toggleProperty('isEditing');
        this.sendAction();
      },

      cancel() {
        this.set('isEditing', false);
      },
    }
});
