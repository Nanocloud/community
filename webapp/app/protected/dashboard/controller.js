import Ember from 'ember';

export default Ember.Controller.extend({

  appData: Ember.inject.service('application-data'),
  session: Ember.inject.service('session'),

  apps: Ember.computed('model.@each', 'model.@each', function() {
    return this.get('model.apps')
  }),

  users: Ember.computed('model', 'model', function() {
    return this.get('model.users')
  }),

  updateSession: function() {

    Ember.$.ajax({
      type: "GET",
      headers: { Authorization : "Bearer " + this.get('session.access_token')},
      url: "api/sessions",
    })
    .then((response) => {
      this.set('sessions', JSON.parse(response).data);
    });

  }.on('init'),

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
