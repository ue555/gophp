<?php
echo "before|";
if (false) { echo "bad"; }
else if (1 / 0) { echo "bad"; }
else { echo $never; }
