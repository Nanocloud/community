import Ember from 'ember';

export default Ember.Controller.extend({
  showSingleTab: false,
  showFileExplorer: false,
  connectionName: null,
  store: Ember.inject.service('store'),
  remoteSession: Ember.inject.service('remote-session'),
  session: Ember.inject.service('session'),

  applicationList: function() {
    this.set('applicationList', this.getFilteredApplicationList());
    return this.getFilteredApplicationList();
  }.property(),

  getFilteredApplicationList: function() {
    return this.get('model').toArray()
      .rejectBy('alias', 'hapticPowershell')
      .rejectBy('alias', 'hapticDesktop');
  },

  actions: {

    updateModel() {
      this.model.update();
    },

    disconnectGuacamole() {
      this.get('remoteSession').disconnectSession(this.get('connectionName'));
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
