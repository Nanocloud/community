import Ember from 'ember';

export default Ember.Controller.extend({

  session: Ember.inject.service('session'),

  loadState: {
    application: 0,
    user: 0,
    session: 0,
  },


  users: Ember.computed('model.users', 'model.users', function() {
    return this.get('model.users')
      .rejectBy('isAdmin', true);
  }),

  apps: Ember.computed('model.apps', 'model.apps', function() {
    return this.get('model.apps')
      .rejectBy('alias', 'hapticDesktop');
  }),

  sessions: Ember.computed('model.@each', 'model.@each', function() {
    return this.get('model.sessions')
      .rejectBy('username', 'Administrator');
  }),

  activator: function() {
    this.loadData('application', 'apps');
    this.loadData('user', 'users');
    this.loadData('session', 'sessions');
  }.on('init'),

  loadData(data, dest) {
    this.set('loadState.' + data, 1);
    this.get('store').query(data, {})
      .then((response) => {
        this.set('loadState.' + data, 0);
        this.set('model.' + dest, response);
      })
      .catch((error) => {
        this.set('loadState.' + data, 2);
      });
  },

  actions : {
    goToApps() {
      this.transitionToRoute('protected.apps');
    },
    goToUsers() {
      this.transitionToRoute('protected.users');
    },
    goToConnectedUsers() {
      this.transitionToRoute('protected.histories');
    },
  }
});
