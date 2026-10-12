<?php
echo "before|";
if (false) { echo "bad"; }
elseif ((+$missing) > (1 / 0)) { echo "bad"; }
else if ($never) { echo "bad"; }
