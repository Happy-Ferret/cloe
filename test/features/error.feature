Feature: Error
  Scenario: Run an erroneous code
    Given a file named "main.tisp" with:
    """
    (write (+ 1 true))
    """
    When I run `tisp main.tisp`
    Then the exit status should not be 0
    And the stdout should contain exactly ""
    And the stderr should contain "Error"
    And the stderr should contain "main.tisp"
    And the stderr should contain "(write (+ 1 true))"

  Scenario: Bind 2 values to an argument
    Given a file named "main.tisp" with:
    """
    (let (f x)
         (let (g y)
              (let (h z) (+ x y z))
              h)
         g)

    (write (((f 123) 456) . x 0))
    """
    When I run `tisp main.tisp`
    Then the exit status should not be 0
    And the stdout should contain exactly ""
    And the stderr should contain "Error"
    And the stderr should contain "main.tisp"
    And the stderr should contain "(write (((f 123) 456) . x 0))"
