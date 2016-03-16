import Ember from 'ember';

export default Ember.Component.extend({
  remoteSession: Ember.inject.service('remote-session'),
  classNames: ['remote-session'],

  guacamole: null,

  connect: function() {

    this.get('guacamole').then((guac) => {
      this.get('element').appendChild(guac.getDisplay().getElement());
      guac.connect();
    });
  }.observes('guacamole').on('didInsertElement')

});
