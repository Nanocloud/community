import Ember from 'ember';

export default Ember.Component.extend({

    originalValue: "",
    isEditing: false,
    errorMessage: null,

    getInputType: function() {
      if (this.get('hideInput')) {
        return "password";
      }
      else {
        return "text";
      }
    }.property('hideInput'),

    isValid: function() {
      return this.get('errorMessage');
    },

    autoSelectInput: function() {
      if (this.get('isEditing')) {
        Ember.run.scheduleOnce('afterRender', () => {
          this.$(this.get('element')).find('input').first().select();
        })
      }
    }.observes('isEditing'),

    actions: {

      toggle() {
        if (this.get('isEditing')) {
          this.set('textInput', this.get('originalValue'));
          this.set('errorMessage', "");
        }

        this.toggleProperty('isEditing');
      },

      submit() {
        var defer = Ember.RSVP.defer();

        defer.promise.then(() => {
          this.set('originalValue', this.get('textInput'));
          this.send('toggle');
        }, (err) => {
          this.set('errorMessage', err);
        });

        this.sendAction('onClose', defer);
      },

      cancel() {
        this.send('toggle');
      },
    }
});
