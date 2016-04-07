import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,
  store: Ember.inject.service('store'),
  remoteSession: Ember.inject.service('remote-session'),
  session: Ember.inject.service('session'),

  applications: Ember.computed(function() {
    return this.get('model')
      .rejectBy('alias', 'hapticPowershell')
      .rejectBy('alias', 'hapticDesktop');
  }),

  actions: {

    disconnectGuacamole(connectionName) {
      this.get('remoteSession').disconnectSession(connectionName);
    },

    publish() {
      this.store.createRecord('application', {});
    },

    toggleSingleTab(connectionName) {
      this.set('connectionName', connectionName);
      this.toggleProperty('showSingleTab');
    },

    toggleFileExplorer() {
      this.toggleProperty('showFileExplorer');
    },
  }
});
