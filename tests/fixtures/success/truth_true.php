<?php
if (true) { echo '1'; } else { echo '0'; }
if (1) { echo '1'; } else { echo '0'; }
if (-1) { echo '1'; } else { echo '0'; }
if (9223372036854775807) { echo '1'; } else { echo '0'; }
if ((-9223372036854775807 - 1)) { echo '1'; } else { echo '0'; }
if (.5) { echo '1'; } else { echo '0'; }
if (-.5) { echo '1'; } else { echo '0'; }
if (5e-324) { echo '1'; } else { echo '0'; }
if ('00') { echo '1'; } else { echo '0'; }
if ('0.0') { echo '1'; } else { echo '0'; }
if ('0e0') { echo '1'; } else { echo '0'; }
if (' 0 ') { echo '1'; } else { echo '0'; }
if (' ') { echo '1'; } else { echo '0'; }
if ('false') { echo '1'; } else { echo '0'; }
if ('null') { echo '1'; } else { echo '0'; }
if ("\0") { echo '1'; } else { echo '0'; }
if ("0\0") { echo '1'; } else { echo '0'; }
if ('word') { echo '1'; } else { echo '0'; }
if ((1e308 * 10)) { echo '1'; } else { echo '0'; }
if (((1e308 * 10) - (1e308 * 10))) { echo '1'; } else { echo '0'; }
echo "\n";
