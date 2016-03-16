import Ember from 'ember';

export default Ember.Component.extend({
  remoteSession: Ember.inject.service('remote-session'),

  guacamole: null,
  connectionName: null,

  connect: function() {

    if (Ember.isEmpty(this.get('connectionName'))) {
      return ;
    }

    let width = $(this.element).parent().width();
    let height = $(this.element).parent().height();

    this.set('guacamole', this.get('remoteSession').getSession(this.get('connectionName'), width, height));
    this.get('guacamole').then((guac) => {
      this.get('element').appendChild(guac.getDisplay().getElement());
      guac.connect();
    });
  }.observes('connectionName').on('becameVisible'),

});
