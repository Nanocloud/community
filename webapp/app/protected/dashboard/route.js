import Ember from 'ember';

export default Ember.Route.extend({

  model() {
    console.log('get get');
    return Ember.RSVP.hash({
      apps: this.store.findAll('application'),
      users: this.store.findAll('user')
    })
  },

});
