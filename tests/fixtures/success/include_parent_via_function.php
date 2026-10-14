<?php
function loadReviewBase() {
    require __DIR__ . '/../support/review_base.inc';
}
loadReviewBase();
class FunctionLoadedChild extends ReviewBase {}
echo (new FunctionLoadedChild)->value(), "\n";
