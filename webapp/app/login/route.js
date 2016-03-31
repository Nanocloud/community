import Ember from 'ember';

export default Ember.Route.extend({
  redirect() {
    if (this.get('session.isAuthenticated')) {
      this.transitionTo('protected.index');
    }
  }
});
