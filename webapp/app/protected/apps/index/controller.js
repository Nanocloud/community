import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,
  store: Ember.inject.service('store'),
  remoteSession: Ember.inject.service('remote-session'),
  session: Ember.inject.service('session'),

  getFileteredApplicationList: function() {
    return this.get('model').toArray()
      .rejectBy('alias', 'hapticPowershell')
      .rejectBy('alias', 'hapticDesktop');
  }.property('model'),

  applications: function() {
    return this.getFileteredApplicationList(this.get('model'));
  }.property(),


  updateApplicationList: function() {
    this.set('applications', this.getFileteredApplicationList());
  },

  actions: {

    updateModel() {
      this.model.update()
    },

    disconnectGuacamole(connectionName) {
      this.get('remoteSession').disconnectSession(connectionName);
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
