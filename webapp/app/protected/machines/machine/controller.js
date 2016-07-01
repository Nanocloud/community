import Ember from 'ember';

export default Ember.Controller.extend({

  controlsSupported: ['qemu','manual'],
  controlsAreSupported: Ember.computed('model.getPlatform', 'controlsSupported', function() {
    var ret = this.get('controlsSupported').indexOf(this.get('model.platform'));
    return ret === -1 ? false : true;
  }),

  machineName: Ember.computed('model.name', function() {
    return this.get('model.name') ? this.get('model.name') : "Machine";
  }),

  startMachine() {
    let machine = this.get('model');

    machine.set('status', 'up');
    machine.save();
  },

  stopMachine() {
    let machine = this.get('model');

    machine.set('status', 'down');
    machine.save();
  },

  terminateMachine() {
    let machine = this.get('model');

    machine.destroyRecord();
    this.transitionToRoute('protected.machines');
  },

  actions: {
    startMachine() {
      this.startMachine();
    },

    stopMachine() {
      this.stopMachine();
    },

    terminateMachine() {
      this.terminateMachine();
    }
  }
});
