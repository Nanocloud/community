import Ember from 'ember';

export default Ember.Component.extend({

    originalValue: "",
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
        if (this.get('isEditing')) {
          this.set('textInput', this.get('originalValue'));
        }

        this.toggleProperty('isEditing');
      },

      submit() {
        var defer = Ember.RSVP.defer();

        defer.promise.then(() => {
          this.set('originalValue', this.get('textInput'));
          this.send('toggle');
        });

        this.sendAction('onClose', defer);
      },

      cancel() {
        this.send('toggle');
      },
    }
});
