import TooltipsterComponent from 'ember-cli-tooltipster/components/tool-tipster';

export default TooltipsterComponent.extend({

  classNameBindings: ['stateEnabled'],

  actions: {
    clickAction() {
      this.sendAction("click");
    }
  }
});
