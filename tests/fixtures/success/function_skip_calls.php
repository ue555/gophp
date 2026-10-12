<?php
function danger() { echo 'bad'; return 1 / 0; }
false && danger(); true || undefinedFunction();
if (false) { danger(); }
while (false) { danger(); }
for (; true; danger()) { break; }
do { break; } while (danger());
echo 'done', "\n";
