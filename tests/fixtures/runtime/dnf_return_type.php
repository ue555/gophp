<?php
interface ReturnA {}
interface ReturnB {}
function invalidReturn(): (ReturnA&ReturnB)|int { return []; }
invalidReturn();
