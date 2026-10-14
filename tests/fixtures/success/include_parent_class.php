<?php
require __DIR__ . '/../support/review_base.inc';
class ReviewChild extends ReviewBase {}
echo (new ReviewChild)->value(), "\n";
