<?php
echo 'before';
class FirstCycle extends SecondCycle {}
class SecondCycle extends FirstCycle {}
echo 'after';
