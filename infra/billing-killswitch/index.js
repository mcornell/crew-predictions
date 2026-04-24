const { CloudBillingClient } = require('@google-cloud/billing');
const billing = new CloudBillingClient();

exports.killBilling = async (pubsubEvent) => {
  const data = JSON.parse(Buffer.from(pubsubEvent.data, 'base64').toString());
  if (data.costAmount <= data.budgetAmount) return;
  const projectIDs = process.env.PROJECT_IDS.split(',');
  for (const projectID of projectIDs) {
    await billing.updateProjectBillingInfo({
      name: `projects/${projectID}`,
      projectBillingInfo: { billingAccountName: '' },
    });
    console.log(`Billing disabled for projects/${projectID}`);
  }
};
