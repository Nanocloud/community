import Ember from 'ember';

export function stateMachine([val]/*, hash*/) {

  var states = {
    unknown: "Unknown",
    down: "Offline",
    up: "Online",
    terminated: "Terminated",
    booting: "Booting in progress",
    downloading: "Downloading",
  };

  return states[val] ? states[val] : val;
}

export default Ember.Helper.helper(stateMachine);
