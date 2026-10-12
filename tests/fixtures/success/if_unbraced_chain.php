<?php
if (true) echo 'A'; elseif (true) echo 'bad'; else echo 'bad';
if (false) echo 'bad'; elseif (true) echo 'B'; else echo 'bad';
if (false) echo 'bad'; else if (false) echo 'bad'; else if (true) echo 'C';
if (false) echo 'bad'; else echo 'D';
if (false) echo 'bad';
if (true); else echo 'bad';
{ echo 'E'; }
echo "\n";
