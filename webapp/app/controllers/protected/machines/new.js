import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    createMachine() {
      let m = this.store.createRecord('machine', {
        name: this.get('machineName'),
        adminPassword: this.get('adminPassword'),
      });

      m.save()
      .then((machine) => {
        this.transitionToRoute('protected.machines.machine', machine);
      });
    }
  }
});
