import Ember from 'ember';
import formatTimeDuration from 'nanocloud/utils/format-duration';

export function formatDuration([value]/*, hash*/) {
  return formatTimeDuration(value);
}

export default Ember.Helper.helper(formatDuration);
