import Ember from 'ember';

export default Ember.Component.extend({


  remoteSession: Ember.inject.service('remote-session'),

  inputFocusChanged: function() {
    if (this.$().find('textarea').length != 0) {
      this.$().find('textarea')
        .focusin(function() {
          this.get('remoteSession').pauseInputs(this.get('connectionName'));
        }.bind(this))
        .focusout(function() {
          this.get('remoteSession').restoreInputs(this.get('connectionName'));
        }.bind(this));
    }
  }.on('didInsertElement'),


  actions: {
    toggleVdiWindow() {
      if (this.get('toggleWindow')) {
        this.toggleWindow();
      }
      else {
        this.toggleProperty('stateVisible');
      }
    },
  }
});
