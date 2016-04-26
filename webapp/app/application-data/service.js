import Ember from 'ember';

export default Ember.Service.extend({

  store: Ember.inject.service('store'),

  /*
  apps: this.get('store').findAll('application', { reload: true }),
  sessions: this.get('store').findAll('session', { reload: true }),
  users: this.get('store').findAll('user', { reload: true }),
  */
//  apps: this.store.all('application')



  init() {
    /*
    this._super(...arguments);
    this.get('store').query('application', {}).then(function(data) {
      this.set('apps', data.rejectBy('alias', 'hapticDesktop'));
    }.bind(this));
    $.getJSON('api/applications/connections', function(data) {
      this.set('sessions', data);
    }.bind(this));
    this.get('store').query('user', {}).then(function(data) {
      this.set('users', data.rejectBy('isAdmin', true));
    }.bind(this));
    */
  },

});
