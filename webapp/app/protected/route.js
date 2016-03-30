import Ember from 'ember';

export default Ember.Route.extend({
  beforeModel() {
    if (this.get('session.isAuthenticated') === false) {
      this.transitionTo('login');
    }
  },

  afterModel(user) {
    this.set('session.user', user);
  },

  model: function() {
    this.set('session.access_token', this.get('session.data.authenticated.access_token'));

    return this.store.queryRecord('user', { me: true })
    .catch((err) => {
      Ember.Logger.error(err);
      this._invalidateSession();
    });
  },

  _invalidateSession: function() {
      this.get('session').invalidate();
      this.transitionTo('login');
  },

  actions: {
    invalidateSession: function() {
      this._invalidateSession();
    }
  }
});
