<?php foreach([0,1,2,3] as $x) { ?>
<?php if($x===0): ?>zero<?php elseif($x===1): ?>one<?php else: ?>other<?php if($x===3): ?>three<?php endif; ?><?php endif; ?>!
<?php } ?>
