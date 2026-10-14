<?php
switch (1) {
    case 1:
        require __DIR__ . '/../support/review_base.inc';
        break;
}
class SwitchLoadedChild extends ReviewBase {}
echo (new SwitchLoadedChild)->value(), "\n";
