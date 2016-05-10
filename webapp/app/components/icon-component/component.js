import TooltipsterComponent from 'ember-cli-tooltipster/components/tool-tipster';

export default TooltipsterComponent.extend({
  classNames: ['icon-component'],
  classNameBindings: [
    'hover-enabled:hover-enabled',
    'clickable:clickable'
  ],
});
