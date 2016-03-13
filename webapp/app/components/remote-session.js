import Ember from 'ember';

export default Ember.Component.extend({
  remoteSession: Ember.inject.service('remote-session'),
  classNames: ['remote-session'],

  connect: function() {
    this.get('remoteSession').set('connectionName', 'hapticDesktop');
    this.get('remoteSession').set('canva', this);
  }.on('didInsertElement')
});
