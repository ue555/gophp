<?php
const X = 3; define("Y", 4); echo X, "|", Y, "|", constant("Y"), "|", defined("X"), "|", defined("missing") ? 1 : 0, "\n";
