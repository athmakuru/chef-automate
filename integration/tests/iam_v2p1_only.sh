#shellcheck disable=SC2034
#shellcheck disable=SC2039
#shellcheck disable=SC2154
#shellcheck disable=SC1091

test_name="iam v2.1 only"
# Note: the inspec profile takes care of deleting any migrated v1 legacy
# policies
test_deploy_inspec_profiles=(a2-iam-v2p1-only-integration)
# Note: we can't run diagnostics AND inspec, so skip diagnostics
test_skip_diagnostics=true

do_setup() {
    do_setup_default

    # We are defaulting to a umask of 077 to test
    # installations on systems that are super locked down.
    # Briefly override that strict default so we can install
    # packages that non-root users can use (like the hab user
    # for health checks or this script).
    local previous_umask
    previous_umask=$(umask)
    umask 022

    hab pkg install core/curl
    hab pkg install -b core/jo

    umask $previous_umask
}

hab_curl() {
    hab pkg exec core/curl curl "$@"
}

do_deploy() {
    do_deploy_default
    log_info "applying dev license"
    chef-automate license apply "$A2_LICENSE"
}

do_test_deploy() {
    log_info "run chef-automate iam upgrade-to-v2 --beta2.1"
    chef-automate iam upgrade-to-v2 --beta2.1 || return 1

    remove_legacy_policies

    do_test_deploy_default
}

remove_legacy_policies() {
    local token=$(chef-automate iam token create ADMIN_TEST --admin)
    local legacy_policies=(secrets-access-legacy
            events-access-legacy
            infrastructure-automation-access-legacy
            compliance-access-legacy
            nodes-access-legacy
            compliance-profile-access-legacy
            ingest-access-legacy
            node-managers-access-legacy
            telemetry-access-legacy)

    for id in "${legacy_policies[@]}"
    do
        echo "Deleting legacy policy $id..."
        curl -sSkH "api-token: $token" -X DELETE https://localhost/apis/iam/v2beta/policies/$id
    done
}
