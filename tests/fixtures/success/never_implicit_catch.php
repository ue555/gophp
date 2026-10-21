<?php
function broken(): never { echo 'inside|'; }
try { broken(); echo 'wrong'; }
catch (TypeError $e) { echo $e->getMessage(), "\n"; }
