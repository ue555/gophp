<?php
if (false) { echo '1'; } else { echo '0'; }
if (null) { echo '1'; } else { echo '0'; }
if (0) { echo '1'; } else { echo '0'; }
if (0.0) { echo '1'; } else { echo '0'; }
if (-0.0) { echo '1'; } else { echo '0'; }
if ('') { echo '1'; } else { echo '0'; }
if ('0') { echo '1'; } else { echo '0'; }
if ((1 - 1)) { echo '1'; } else { echo '0'; }
if (('1' - 1)) { echo '1'; } else { echo '0'; }
if ((1 === 2)) { echo '1'; } else { echo '0'; }
if (('0' . '')) { echo '1'; } else { echo '0'; }
if (1e-9999) { echo '1'; } else { echo '0'; }
if (-1e-9999) { echo '1'; } else { echo '0'; }
echo "\n";
