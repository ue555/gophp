<?php
echo "before|";
require __DIR__ . '/../support/review_base.inc';
class TimedChild extends ReviewBase {}
echo (new TimedChild)->value(), "\n";
