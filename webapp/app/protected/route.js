import Ember from 'ember';

export default Ember.Route.extend({

  beforeModel(transition) {
    if (transition.queryParams.app) {  
      this.set('directLinkParams', transition.queryParams);
    }
    if (this.get('session.isAuthenticated') === false) {
      this.transitionTo('login');
    }
    else {
      if (this.get('directLinkParams')) {
        this.transitionTo('direct-link', {
          queryParams: this.get('directLinkParams')
        });
      }
    }
  },

  redirect() {

    if (this.get('session.isAuthenticated') === true) {

      if (this.get('session.user.isAdmin')) {
        this.transitionTo('protected.dashboard');
      }
      else {
        this.transitionTo('protected.apps');
      }
    }
  },

  afterModel(user) {
    this.set('session.user', user);
  },

  model() {
    this.set('session.access_token', this.get('session.data.authenticated.access_token'));
    return this.store.queryRecord('user', { me: true });
  }
});
