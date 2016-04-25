import Ember from 'ember';

export default Ember.Component.extend({


  remoteSession: Ember.inject.service('remote-session'),
  hasFocus: false,

  mouseEnter() {
    console.log("ENTER");
    this.set('hasFocus', false);
    this.get('remoteSession').pauseInputs(this.get('connectionName'));
  },

  mouseLeave() {
    console.log("LEAVE");
    this.set('hasFocus', true);
    this.get('remoteSession').restoreInputs(this.get('connectionName'));
  },

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
