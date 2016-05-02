import Ember from 'ember';

export default Ember.Controller.extend({
  reset: function() {
    this.setProperties({
      machineName: '',
      adminPassword: ''
    });
  },

  actions: {
    createMachine() {
      let type = this.get('types').objectAt(0);

      let m = this.store.createRecord('machine', {
        name: this.get('machineName'),
        adminPassword: this.get('adminPassword'),
        type: type
      });

      m.save()
      .then((machine) => {
        this.transitionToRoute('protected.machines.machine', machine);
      });
    }
  }
});
