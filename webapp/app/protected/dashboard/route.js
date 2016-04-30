import Ember from 'ember';

export default Ember.Route.extend({

  loading: function() {
    var view = this.container.lookup('view:loading').append();
    this.router.one('didTransition', view, 'destroy');
  },

  model() {
    return Ember.RSVP.hash({
      apps: this.store.findAll('application'),
      users: this.store.findAll('user'),
      sessions: this.store.findAll('session')
    });
  },
});
