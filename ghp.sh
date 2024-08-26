#/bin/bash

if oco; then
    echo "Eco command succeeded, proceeding with GitHub checks..."

    # Capture the PR state using gh CLI
    pr_state=$(gh pr view --json state --jq '.state' 2>/dev/null)

    # Check the exit status of the gh pr view command and the PR state
    if [[ $? -eq 0 && "$pr_state" == "OPEN" ]]; then
        # If the exit status is 0 and the PR state is OPEN, run the following command
        gh pr create --web
    else
        # If the exit status is not 0 or the PR state is not OPEN, do nothing
        echo "Pull request is not open, or an error occurred. No action taken."
    fi
else
    # If the eco command fails, do nothing or handle the error
    echo "Eco command failed. No further actions taken."
fi
